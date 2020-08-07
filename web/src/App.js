import React from 'react';
import './App.css';
import { Route, BrowserRouter, Switch } from 'react-router-dom';
import 'bootstrap/dist/css/bootstrap.min.css';
import { loadStripe } from "@stripe/stripe-js";
import { Elements } from "@stripe/react-stripe-js";

// Layout
import Navbar from './components/layout/Navbar'

// Pages
import Main from './components/pages/Main'
import Users from './components/pages/Users'
import Products from './components/pages/Products'
import Shops from './components/pages/Shops'
import Reviews from './components/pages/Reviews'
import Payment from './components/pages/Payment'

const promise = loadStripe("pk_test_51HDBdpF2hb5ZuDX8NV3QCssBlpln8sHMwxA3CwgYEtzJXbLbDa6BHEbN0NaD1NxR3V79ay6MZVctdBy93okmH1sn00afdDfoBw");

function App() {
  return (
    <BrowserRouter>

    <Navbar />

      <Switch>
        <Route exact path="/" component={Main} />

        <Route exact path="/users" component={Users} />

        <Elements stripe={promise}>
          <Route exact path="/payment" component={Payment} />
        </Elements>

        <Route exact path="/products" component={Products} />

        <Route exact path="/reviews" component={Reviews} />

        <Route exact path="/shops" component={Shops} />

      </Switch>

    </BrowserRouter>
  );
}

export default App;
