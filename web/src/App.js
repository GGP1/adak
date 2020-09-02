import React from 'react';
import './App.css';
import { Route, BrowserRouter, Switch } from 'react-router-dom';
import 'bootstrap/dist/css/bootstrap.min.css';

// Layout
import Navbar from './components/layout/Navbar';

// Pages
import Login from './components/pages/Login';
import Main from './components/pages/Main';
import Payment from './components/pages/Payment';
import Products from './components/pages/Products';
import Register from './components/pages/Register';
import Shops from './components/pages/Shops';
import Users from './components/pages/Users';


function App() {
  return (
    <BrowserRouter>

    <Navbar />

      <Switch>
        <Route exact path="/" component={Main} />

        <Route exact path="/login" component={Login} />

        <Route exact path="/payment" component={Payment} /> {/* not working */}

        <Route exact path="/products" component={Products} />

        <Route exact path="/register" component={Register} />

        <Route exact path="/shops" component={Shops} />

        <Route exact path="/users" component={Users} />
      </Switch>

    </BrowserRouter>
  );
}

export default App;
