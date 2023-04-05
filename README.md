# zerok CLI

## Install zkcli

### From Script

CLI now has an installer script that will automatically grab the latest version and install it locally.

```sh
bash -c "$(curl -fsSL https://zerok.ai/install.sh)"
```

### From Script

CLI now has an installer script that will automatically grab the latest version and install it locally.

### From the Binary Releases

Binary downloads of the CLI can be found on [the Releases page](https://https://github.com/zerok-ai/zk-cli/releases/latest).


## Install zerok

Run the following command to install zerok in the current cluster context

```sh
zkctl install --apikey [api-key]
```

The `api-key` is available through our [dashboard](http://dashboard.zerok.ai/api-key).

## Activate zerok and do rolling restart

Each namespace in the cluster has to be marked for ZeroK. Once marked, all the new pods will get activated for zerok. I can be done using the following command:

```sh
zkctl activate -n <namespace>
```

You have to restart the old pods. You can do both activation and restart using the following command:

```sh
zkctl activate -n <namespace> -r
```
