import React from 'react';
import { ApolloProvider } from '@apollo/client/react';
import { client } from './lib/apolloClient';
import { Component } from './Component';

function App() {
  return (
    <ApolloProvider client={client}>
      <Component />
    </ApolloProvider>
  );
}

export default App;
