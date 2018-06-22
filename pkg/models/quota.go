// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"fmt"
)

type Quota struct {
	Name  string
	Count int
}

type Quotas struct {
	Instance   *Quota
	Cpu        *Quota
	Gpu        *Quota
	Memory     *Quota
	Volume     *Quota
	VolumeSize *Quota
}

func NewQuotas() *Quotas {
	quotas := &Quotas{
		Instance:   new(Quota),
		Cpu:        new(Quota),
		Gpu:        new(Quota),
		Memory:     new(Quota),
		Volume:     new(Quota),
		VolumeSize: new(Quota),
	}
	return quotas
}

func (p *Quotas) LessThan(quotas *Quotas) error {
	if p.Instance.Count > quotas.Instance.Count {
		return fmt.Errorf("need %d more %s quota", p.Instance.Count-quotas.Instance.Count, p.Instance.Name)
	}
	if p.Cpu.Count > quotas.Cpu.Count {
		return fmt.Errorf("need %d more %s quota", p.Cpu.Count-quotas.Cpu.Count, p.Cpu.Name)
	}
	if p.Gpu.Count > quotas.Gpu.Count {
		return fmt.Errorf("need %d more %s quota", p.Gpu.Count-quotas.Gpu.Count, p.Gpu.Name)
	}
	if p.Memory.Count > quotas.Memory.Count {
		return fmt.Errorf("need %d more %s quota", p.Memory.Count-quotas.Memory.Count, p.Memory.Name)
	}
	if p.Volume.Count > quotas.Volume.Count {
		return fmt.Errorf("need %d more %s quota", p.Volume.Count-quotas.Volume.Count, p.Volume.Name)
	}
	if p.VolumeSize.Count > quotas.VolumeSize.Count {
		return fmt.Errorf("need %d more %s quota", p.VolumeSize.Count-quotas.VolumeSize.Count, p.VolumeSize.Name)
	}
	return nil
}
