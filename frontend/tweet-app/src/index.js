import React from 'react';
import { ApolloClient, ApolloProvider, InMemoryCache } from '@apollo/client';
import { createRoot } from 'react-dom/client';  // Import createRoot from "react-dom/client"
import {App} from './App';
import reportWebVitals from './reportWebVitals';

const client = new ApolloClient({
  url: 'http://localhost:8080/graphql',
  cache: new InMemoryCache(),
});

const root = createRoot(document.getElementById('root'));

root.render(
  <React.StrictMode>
    <ApolloProvider client={client}>
      <App />
    </ApolloProvider>
  </React.StrictMode>,
);

reportWebVitals();