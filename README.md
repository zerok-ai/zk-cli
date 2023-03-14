# Install ZeroK

## Add alias (for development)

If you are running the code instead of the executable set the following alias

``` sh
alias zkctl="go run main.go"
```

## download zerok CLI

Copy and run the command mentioned below to install the zerok CLI.

```sh
bash -c "$(curl -fsSL https://zerok.ai/install.sh)"
```

## install zerok

Run the following command to install zerok on the current cluster context

```sh
zkctl install 
```

## activate zerok and do rolling restart

Each namespace in the cluster has to be marked for ZeroK. Once marked, all the new pods will get activated for zerok. I can be done using the following command:

```sh
zkctl activate -n <namespace>
```

You have to restart the old pods. You can do both activation and restart using the following command:

```sh
zkctl activate -n <namespace> -r
```