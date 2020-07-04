import React from 'react';
import './App.css';
import { Route, BrowserRouter, Switch } from 'react-router-dom';
import 'bootstrap/dist/css/bootstrap.min.css';


// Pages
import Main from './components/pages/Main'
import Users from './components/pages/Users'
import Products from './components/pages/Products'
import Shops from './components/pages/Shops'
import Reviews from './components/pages/Reviews'

function App() {
  return (
    <BrowserRouter>

      <Switch>
        <Route exact path="/" component={Main} />

        <Route exact path="/users" component={Users} />

        <Route exact path="/products" component={Products} />

        <Route exact path="/shops" component={Shops} />

        <Route exact path="/reviews" component={Reviews} />
      </Switch>

    </BrowserRouter>
  );
}

export default App;
