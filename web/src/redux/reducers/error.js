const defaultState = {};

const errorReducer = (state = defaultState, action) => {
  switch (action.type) {
    case 'ERROR':
      return action.payload;
    default:
      return state;
  }
}

export default errorReducer;