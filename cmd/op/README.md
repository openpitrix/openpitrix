# `op`

`op` is a command line tool for easy to create OpenPitrix VM-Based Application.

## How to Use?

```bash
# create your repo directory
mkdir myrepo && cd myrepo
# create your first OpenPitrix App, named with `nginx`
op create nginx
# let's take a look with `nginx` directory
ls nginx
# output: cluster.json.tmpl  config.json  package.json
# ---
# after your edit files under `nginx` directory
# package `nginx` to a archived file
op package nginx
# output: Successfully packaged chart and saved it to: /$YOURPATH/myrepo/chart/nginx-0.1.0.tgz

# generte index.yaml, so that the repo_indexer of OpenPitrix can create app with this repo
op index ./ && cat index.yaml
# your can serve a local repo now
op serve .
# if you want to publish the repo, regenerate index.yaml with your website link `http://example.cn/`
op index --url http://example.cn/ ./
# you can upload `index.yaml` and `*.tgz` to your site now
```

