import React, { useEffect, useState } from 'react';
import { ApolloProvider } from '@apollo/client/react';
import { client } from './lib/apolloClient';
import { Component } from './Component';
import { NotFound } from './NotFound';

function App() {
  const [currentPath, setCurrentPath] = useState(window.location.pathname);

  useEffect(() => {
    const handleNavigation = () => {
      setCurrentPath(window.location.pathname);
    };

    window.addEventListener('popstate', handleNavigation);
    return () => window.removeEventListener('popstate', handleNavigation);
  }, []);

  const renderContent = () => {
    if (currentPath === '/') {
      return <Component />;
    }
    return <NotFound />;
  };

  return (
    <ApolloProvider client={client}>
      {renderContent()}
    </ApolloProvider>
  );
}

export default App;
