import React, { useState, useEffect } from 'react';
import axios from 'axios';

import { IterateShops } from '../subcomponents/Iterations/Iterations'

function Shops() {
    const [shops, setShops] = useState([])

    useEffect(() => {
        async function getShops() {
            const res = await axios.get('http://localhost:4000/shops')
            // use custom headers to fetch the cookie value from the server and
            // set front-end cookies with the same values
            // then, when requesting something that requires being logged in, send the cookie
            // value in another header to the server
            console.log(res.headers.authorization); // undefined if there's no header set
            setShops(res.data)
        }
        getShops();
    }, [])

    return (
        <div className="container">
            <IterateShops shops={shops} />
        </div>
    )
}

export default Shops;