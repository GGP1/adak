import React from 'react';
import './App.css';
import { Route, BrowserRouter, Switch } from 'react-router-dom';
import 'bootstrap/dist/css/bootstrap.min.css';

// Layout
import LeftBar from './components/layout/LeftBar'

// Pages
import Main from './components/pages/Main'

function App() {
  return (
    <BrowserRouter>

      <LeftBar />

      <Switch>
        <Route exact path="/" component={Main} />
      </Switch>

    </BrowserRouter>
  );
}

export default App;
