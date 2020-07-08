import { LOGIN, USER_LOADING, LOGOUT } from '../actions/types'

const defaultState = {
  authenticated: false,
  token: null,
  loading: false
};

const authReducer = (state = defaultState, action) => {
  switch (action.type) {
    case LOGIN:
      return {
        ...state,
        authenticated: true,
        token: action.payload
      };
      case LOGOUT:
        return {
          ...state,
          authenticated: false
        };
    case USER_LOADING:
      return {
        ...state,
        loading: true
      };
    default:
      return state;
  }
}

export default authReducer;