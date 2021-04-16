# Branching plugin for Traefik

This plugin enables the use of conditional middleware chains in a router.

Altough all the same middleware from Traefik are available, since we don't have access to the runtime environment we can't just reference existing middleware, instead a new chain is build based on the plugin configuration.

This is just an experiment and should be better implemented as a standard Traefik middleware with access to the runtime.

# Config Example

```yaml
Condition: "Header[`Foo`].0 == `bar`"
  Chain:
    test-prefix:
      AddPrefix:
        Prefix: "/test"
```

On this example only when the request header `Foo` has value `bar` the prefix `/test` is added before sending to the backend server.