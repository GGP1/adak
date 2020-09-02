import React from "react";

import { CardElement } from "@stripe/react-stripe-js";

function Payment() {
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

  return (
    <form style={formStyle} id="payment-form">
      <CardElement id="card-element" options={cardStyle} />
    </form>
  );
}

export default Payment;