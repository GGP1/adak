import React, { useState, useEffect } from "react";
import axios from 'axios';
import Cookies from 'js-cookie';

import {
  CardElement,
  useStripe,
  useElements
} from "@stripe/react-stripe-js";

export default function CheckoutForm() {
  const [succeeded, setSucceeded] = useState(false);
  const [error, setError] = useState(null);
  const [processing, setProcessing] = useState('');
  const [disabled, setDisabled] = useState(true);
  const [clientSecret, setClientSecret] = useState('');
  const [email, setEmail] = useState('');

  const stripe = useStripe();
  const elements = useElements();

  const cardStyle = {
    style: {
      base: {
        color: "#32325d",
        fontFamily: 'Arial, sans-serif',
        fontSmoothing: "antialiased",
        fontSize: "16px",
        "::placeholder": {
          color: "#32325d"
        }
      },
      invalid: {
        color: "#fa755a",
        iconColor: "#fa755a"
      }
    }
  };

  const formStyle = {
    alignItems: 'center',
    width: '45vw',
    margin: 'auto',
    marginTop: '50px'
  }

  const handleChange = async (event) => {
    // Listen for changes in the CardElement
    // and display any errors as the customer types their card details
    setDisabled(event.empty);
    setError(event.error ? event.error.message : "");
  };

  // On submit, send user card details to the server.
  const handleSubmit = async ev => {
    ev.preventDefault();
    setProcessing(true);
  };

  return (
    <form style={formStyle} id="payment-form" onSubmit={handleSubmit}>

      <CardElement id="card-element" options={cardStyle} onChange={handleChange} />

    </form>
  );
}