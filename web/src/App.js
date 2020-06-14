import React from 'react';
import './App.css';
import { Route, BrowserRouter, Switch } from 'react-router-dom';
import 'bootstrap/dist/css/bootstrap.min.css';

// Layout
import LeftBar from './components/layout/LeftBar'

// Pages
import Main from './components/pages/Main'
import Users from './components/pages/Users'
import Products from './components/pages/Products'
import Shops from './components/pages/Shops'

function App() {
  return (
    <BrowserRouter>

      <LeftBar />

      <Switch>
        <Route exact path="/" component={Main} />

        <Route exact path="/users" component={Users} />

        <Route exact path="/products" component={Products} />

        <Route exact path="/shops" component={Shops} />
      </Switch>

    </BrowserRouter>
  );
}

export default App;
