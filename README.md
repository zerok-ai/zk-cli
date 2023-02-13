# Install ZeroK

## Add alias

If you are running the code instead of the executable set the following alias

``` sh
alias zkctl="go run main.go"
```

## install zerok operator

```sh
zkctl install operator
```

## install backend

```sh
zkctl install backend
```

## perform post backend installation tasks

```sh
zkctl install backend postsetup
```

## activate zerok and do rolling restart

Activate a namespace for ZeroK and do a rolling restart

```sh
zkctl activate -n <namespace> -r
```