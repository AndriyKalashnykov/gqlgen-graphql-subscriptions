# front-end

## Run front-end

The following will re-install the node modules and start front-end

```bash
cd frontend

yarn add @apollo/client graphql subscriptions-transport-ws svgo
yarn add @chakra-ui/react @emotion/react @emotion/styled framer-motion
yarn add @chakra-ui/styled-system @babel/core @testing-library/dom @babel/plugin-syntax-flow @babel/plugin-transform-react-jsx
yarn install
yarn upgrade --latest
yarn build
yarn start
```

### Open a second browser

This should open a second browser at [http://localhost:3000](http://localhost:3000).
Post a second message and see it appear in the other window.

Once the application is running, you will notice redis message [XADD](https://redis.io/commands/XADD)
and [XREAD](https://redis.io/commands/xread) appear in redis output.

Note: `XADD` has specified `MAXLEN` argument if 10 to limit the total size or history of the stream,
so if you opened a new window, you would only see up to 10 messages in the history.
