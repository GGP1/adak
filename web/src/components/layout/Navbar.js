import React, { useState } from 'react'

// styles
import '../../App.css';

import SignInButton from './SignInButton';
import DeployableNavbar from './DeployableNavbar';

function Navbar() {
    const [display, setDisplay] = useState(false);

    const toggleNavbar = () => {
        setDisplay(!display)
    }

    return (
        <nav className="navbar navbar-expand-lg navbar-dark" id="nav" style={{ backgroundColor: 'rgba(32, 49, 59, 1)' }}>
            <div className="container mx-auto">
                <a className="navbar-brand text-left" id="brand-text" href="/">
                    Palo
                    </a>
                <button
                    className="navbar-toggler"
                    id="nav-btn"
                    type="button"
                    data-toggle="collapse"
                    data-target="#navbarNav"
                    aria-controls="navbarNav"
                    aria-expanded="false"
                    aria-label="Toggle navigation"
                    onClick={toggleNavbar}
                >
                    <span className="navbar-toggler-icon"></span>
                </button>
                <div className="collapse navbar-collapse" id="navbarNav">
                    <ul className="navbar-nav ml-auto">
                        <li className="nav-item ml-2 px-1 rounded" id="navLi">
                            <a href="/" className="nav-link text-center rounded" id="navLink" style={{ fontFamily: 'Open Sans', color: '#fff' }}>
                                Home
                                </a>
                        </li>
                        <li className="nav-item ml-2 px-1" id="navLi">
                            <a href="/products" className="nav-link text-center rounded" id="navLink" style={{ fontFamily: 'Open Sans', color: 'rgba(255, 255, 255, 0.7)' }}>
                                Products
                                </a>
                        </li>
                        <li className="nav-item ml-2 px-1" id="navLi">
                            <a className="nav-link text-center rounded" href="/shops" id="navLink" style={{ fontFamily: 'Open Sans', color: 'rgba(255, 255, 255, 0.7)' }}>
                                Shops
                                </a>
                        </li>
                        <li className="nav-item ml-2 px-1" style={{ marginTop: '2px' }}>
                            <SignInButton />
                        </li>
                    </ul>
                </div>
            </div>

            {display ?
                <DeployableNavbar />
                : null}

        </nav>
    )

}

export default Navbar;