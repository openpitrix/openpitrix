// Copyright confd. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE-confd file.

package libconfd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"text/template"
)

type TemplateResourceProcessor struct {
	TemplateResource

	path          string
	client        BackendClient
	store         *KVStore
	stageFile     *os.File
	templateFunc  *TemplateFunc
	funcMap       template.FuncMap
	keepStageFile bool
	lastIndex     uint64
	syncOnly      bool
	noop          bool
}

func MakeAllTemplateResourceProcessor(
	config *Config, client BackendClient,
) (
	[]*TemplateResourceProcessor,
	error,
) {
	GetLogger().Debug("Loading template resources from confdir " + config.ConfDir)

	tcs, paths, err := ListTemplateResource(config.GetConfigDir())
	if err != nil {
		if len(paths) == 0 {
			GetLogger().Warning("Found no templates")
			return nil, fmt.Errorf("Found no templates")
		} else {
			GetLogger().Warning(err) // skip error
		}
	}

	templates := make([]*TemplateResourceProcessor, len(paths))
	for i, p := range paths {
		templates[i] = NewTemplateResourceProcessor(
			p, config, client, tcs[i],
		)
	}

	return templates, nil
}

// NewTemplateResourceProcessor creates a NewTemplateResourceProcessor.
func NewTemplateResourceProcessor(
	path string, config *Config, client BackendClient, res *TemplateResource,
) *TemplateResourceProcessor {
	GetLogger().Debug("Loading template resource from " + path)

	tr := TemplateResourceProcessor{
		TemplateResource: *res,
	}

	tr.path = path
	tr.client = client
	tr.store = NewKVStore()
	tr.keepStageFile = config.KeepStageFile
	tr.syncOnly = config.SyncOnly
	tr.noop = config.Noop

	// replace ${LIBCONFD_CONFDIR}
	tr.Dest = strings.Replace(tr.Dest, `${LIBCONFD_CONFDIR}`, config.ConfDir, -1)
	tr.CheckCmd = strings.Replace(tr.CheckCmd, `${LIBCONFD_CONFDIR}`, config.ConfDir, -1)
	tr.ReloadCmd = strings.Replace(tr.ReloadCmd, `${LIBCONFD_CONFDIR}`, config.ConfDir, -1)

	if config.ConfDir != "" {
		if s := tr.Dest; !filepath.IsAbs(s) {
			os.MkdirAll(config.GetDefaultTemplateOutputDir(), 0744)
			tr.Dest = filepath.Join(config.GetDefaultTemplateOutputDir(), s)
			tr.Dest = filepath.Clean(tr.Dest)
		}
	}

	if config.Prefix != "" {
		tr.Prefix = config.Prefix
	}

	if !strings.HasPrefix(tr.Prefix, "/") {
		tr.Prefix = "/" + tr.Prefix
	}

	if len(config.PGPPrivateKey) > 0 {
		tr.PGPPrivateKey = append([]byte{}, config.PGPPrivateKey...)
	}

	if tr.Uid == -1 {
		tr.Uid = os.Geteuid()
	}

	if tr.Gid == -1 {
		tr.Gid = os.Getegid()
	}

	tr.templateFunc = NewTemplateFunc(tr.store, tr.PGPPrivateKey, config.HookAbsKeyAdjuster)
	tr.funcMap = tr.templateFunc.FuncMap

	if !filepath.IsAbs(tr.Src) {
		tr.Src = filepath.Join(config.GetTemplateDir(), tr.Src)
	}

	return &tr
}

// process is a convenience function that wraps calls to the three main tasks
// required to keep local configuration files in sync. First we gather vars
// from the store, then we stage a candidate configuration file, and finally sync
// things up.
// It returns an error if any.
func (p *TemplateResourceProcessor) Process(call *Call) (err error) {
	if fn := call.Config.HookOnUpdateDone; fn != nil {
		defer func() { fn(p.path, err) }()
	}

	if len(call.Config.FuncMap) > 0 {
		for k, fn := range call.Config.FuncMap {
			p.funcMap[k] = fn
		}
	}
	if fn := call.Config.FuncMapUpdater; fn != nil {
		fn(p.funcMap, p.templateFunc)
	}

	if err := p.setFileMode(call); err != nil {
		GetLogger().Error(err)
		return err
	}
	if err := p.setVars(call); err != nil {
		GetLogger().Error(err)
		return err
	}
	if err := p.createStageFile(call); err != nil {
		GetLogger().Error(err)
		return err
	}
	if err := p.sync(call); err != nil {
		GetLogger().Error(err)
		return err
	}
	return nil
}

// setFileMode sets the FileMode.
func (p *TemplateResourceProcessor) setFileMode(call *Call) error {
	if p.Mode == "" {
		if fi, err := os.Stat(p.Dest); err == nil {
			p.FileMode = fi.Mode()
		} else {
			p.FileMode = 0644
		}
	} else {
		mode, err := strconv.ParseUint(p.Mode, 0, 32)
		if err != nil {
			return err
		}
		p.FileMode = os.FileMode(mode)
	}
	return nil
}

// setVars sets the Vars for template resource.
func (p *TemplateResourceProcessor) setVars(call *Call) error {
	GetLogger().Debugln("prefix:", p.Prefix)

	absKeys := p.getAbsKeys()
	GetLogger().Debugf("absKeys: %#v\n", absKeys)

	GetLogger().Debugf("GetValues: absKeys0 = %#v\n", absKeys)

	if fn := call.Config.HookAbsKeyAdjuster; fn != nil {
		for i, key := range absKeys {
			absKeys[i] = fn(key)
		}
	}

	GetLogger().Debugf("GetValues: absKeys1 = %#v\n", absKeys)

	values, err := p.client.GetValues(absKeys)
	if err != nil {
		return err
	}

	GetLogger().Debugf("GetValues: %#v\n", values)

	p.store.Purge()
	for k, v := range values {
		//p.store.Set(path.Join("/", strings.TrimPrefix(k, p.Prefix)), v)
		p.store.Set(k, v)
	}

	return nil
}

// createStageFile stages the src configuration file by processing the src
// template and setting the desired owner, group, and mode. It also sets the
// StageFile for the template resource.
// It returns an error if any.
func (p *TemplateResourceProcessor) createStageFile(call *Call) error {
	if fileNotExists(p.Src) {
		err := errors.New("Missing template: " + p.Src)
		GetLogger().Error(err)
		return err
	}

	tmpl, err := template.New(filepath.Base(p.Src)).Funcs(template.FuncMap(p.funcMap)).ParseFiles(p.Src)
	if err != nil {
		err := fmt.Errorf("Unable to process template %s, %s", p.Src, err)
		GetLogger().Error(err)
		return err
	}

	// create TempFile in Dest directory to avoid cross-filesystem issues
	ensureFileDir(p.Dest)
	temp, err := ioutil.TempFile(filepath.Dir(p.Dest), "."+filepath.Base(p.Dest))
	if err != nil {
		GetLogger().Error(err)
		return err
	}

	if err = tmpl.Execute(temp, nil); err != nil {
		temp.Close()
		os.Remove(temp.Name())
		GetLogger().Error(err)
		return err
	}
	defer temp.Close()

	// Set the owner, group, and mode on the stage file now to make it easier to
	// compare against the destination configuration file later.
	os.Chmod(temp.Name(), p.FileMode)
	os.Chown(temp.Name(), p.Uid, p.Gid)

	p.stageFile = temp
	return nil
}

// sync compares the staged and dest config files and attempts to sync them
// if they differ. sync will run a config check command if set before
// overwriting the target config file. Finally, sync will run a reload command
// if set to have the application or service pick up the changes.
// It returns an error if any.
func (p *TemplateResourceProcessor) sync(call *Call) error {
	staged := p.stageFile.Name()

	if p.keepStageFile {
		GetLogger().Info("Keeping staged file: " + staged)
	} else {
		defer os.Remove(staged)
	}

	GetLogger().Debug("Comparing candidate config to " + p.Dest)

	isSame, err := p.checkSameConfig(staged, p.Dest)
	if err != nil {
		GetLogger().Warning(err)
		return err
	}

	if p.noop {
		GetLogger().Warning("Noop mode enabled. " + p.Dest + " will not be modified")
		return nil
	}
	if isSame {
		GetLogger().Debug("Target config " + p.Dest + " in sync")
		return nil
	}

	GetLogger().Info("Target config " + p.Dest + " out of sync")
	if !p.syncOnly && strings.TrimSpace(p.CheckCmd) != "" {
		if err := p.doCheckCmd(call); err != nil {
			return fmt.Errorf("Config check failed: %v", err)
		}
	}

	GetLogger().Debug("Overwriting target config " + p.Dest)

	err = os.Rename(staged, p.Dest)
	if err != nil {
		GetLogger().Debug("Rename failed - target is likely a mount. Trying to write instead")

		if !strings.Contains(err.Error(), "device or resource busy") {
			return err
		}

		// try to open the file and write to it

		var contents []byte
		var rerr error
		contents, rerr = ioutil.ReadFile(staged)
		if rerr != nil {
			return rerr
		}

		err := ioutil.WriteFile(p.Dest, contents, p.FileMode)
		// make sure owner and group match the temp file, in case the file was created with WriteFile
		os.Chown(p.Dest, p.Uid, p.Gid)
		if err != nil {
			return err
		}
	}

	if !p.syncOnly && strings.TrimSpace(p.ReloadCmd) != "" {
		if err := p.doReloadCmd(call); err != nil {
			return err
		}
	}

	GetLogger().Info("Target config " + p.Dest + " has been updated")
	return nil
}

// check executes the check command to validate the staged config file. The
// command is modified so that any references to src template are substituted
// with a string representing the full path of the staged file. This allows the
// check to be run on the staged file before overwriting the destination config
// file.
// It returns nil if the check command returns 0 and there are no other errors.
func (p *TemplateResourceProcessor) doCheckCmd(call *Call) (err error) {
	if fn := call.Config.HookOnCheckCmdDone; fn != nil {
		defer func() { fn(p.path, p.CheckCmd, err) }()
	}

	var cmdBuffer bytes.Buffer
	data := make(map[string]string)
	data["src"] = p.stageFile.Name()
	tmpl, err := template.New("checkcmd").Parse(p.CheckCmd)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(&cmdBuffer, data); err != nil {
		return err
	}
	return p.runCommand(cmdBuffer.String())
}

// reload executes the reload command.
// It returns nil if the reload command returns 0.
func (p *TemplateResourceProcessor) doReloadCmd(call *Call) (err error) {
	if fn := call.Config.HookOnReloadCmdDone; fn != nil {
		defer func() { fn(p.path, p.ReloadCmd, err) }()
	}

	return p.runCommand(p.ReloadCmd)
}

// runCommand is a shared function used by check and reload
// to run the given command and log its output.
// It returns nil if the given cmd returns 0.
// The command can be run on unix and windows.
func (_ *TemplateResourceProcessor) runCommand(cmd string) error {
	cmd = strings.TrimSpace(cmd)

	GetLogger().Debug("TemplateResourceProcessor.runCommand: " + cmd)

	if _LIBCONFD_GOOS != runtime.GOOS {
		err := fmt.Errorf("cross GOOS(%s) donot support runCommand!", _LIBCONFD_GOOS)
		GetLogger().Error(err)
		return err
	}

	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		c = exec.Command("cmd", "/C", cmd)
	} else {
		c = exec.Command("/bin/sh", "-c", cmd)
	}

	output, err := c.CombinedOutput()
	if err != nil {
		GetLogger().Errorf("%v, output: %q", err, string(output))
		return err
	}

	GetLogger().Debugf("%q", string(output))
	return nil
}

// checkSameConfig reports whether src and dest config files are equal.
// Two config files are equal when they have the same file contents and
// Unix permissions. The owner, group, and mode must match.
// It return false in other cases.
func (_ *TemplateResourceProcessor) checkSameConfig(src, dest string) (bool, error) {
	d, err := readFileStat(dest)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	s, err := readFileStat(src)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return d == s, nil
}
