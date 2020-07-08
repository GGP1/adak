import axios from 'axios';
import setAuthToken from "../utils/setAuthToken";
import jwt_decode from "jwt-decode";
import { LOGIN, USER_LOADING, ERROR, LOGOUT } from './types'

// --- Register ---
export const registerUser = (userData, history) => dispatch => {
    axios
        .post("http://localhost:4000/users", userData)
        .then(res => history.push("/login"))
        .catch(err =>
            dispatch({
                type: ERROR,
                payload: err.response.data
            })
        );
};

// --- Login ---
export const loginUser = userData => dispatch => {
    axios.post('http://localhost:4000/users/login', userData)
        .then(res => {
            // Set token to localStorage
            const { token } = res.data;
            const { message } = res.data;

            localStorage.setItem("jwtToken", token);
            // Set token to Auth header
            setAuthToken(token);
            // Decode token to get user data
            const decoded = jwt_decode(token);
            // Set current user
            dispatch(setCurrentUser(decoded));
            console.log(message);
        })
        .catch(err =>
            dispatch({
                type: ERROR,
                payload: err.response.data
            })
        );
};

// Set logged in user
export const setCurrentUser = decoded => {
    return {
        type: LOGIN,
        payload: decoded
    };
};

// User loading
export const setUserLoading = () => {
    return {
        type: USER_LOADING
    };
};

// Set authenticated to
export const setLogout = () => {
    return {
        type: LOGOUT
    };
};

// --- Logout ---
export const logoutUser = () => dispatch => {
    // Remove token from local storage
    localStorage.removeItem("jwtToken");
    // Remove auth header for future requests
    setAuthToken(false);
    // Set current user to empty object {}
    dispatch(setCurrentUser({}));
    // Set authenticated to false
    dispatch(setLogout());
};