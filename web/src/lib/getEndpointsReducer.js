export const initialStateGetAllEndpoints = {
    endpoints: [],
    loading: true,
    errorMessage: null,
};
export const GetEndpointsReducer = (initialState, action) => {

    switch (action.type) {
        case 'ENDPOINTS_REQUEST':
            return {
                ...initialState,
                loading: true,
            };
        case 'ENDPOINTS_SUCCESS':
            return {
                ...initialState,
                polls: action.payload,
                loading: false,
            };

        case 'ENDPOINTS_ERROR':
            return {
                ...initialState,
                loading: false,
                errorMessage: action.error,
            };

        default:
            throw new Error(`Unhandled action type: ${action.type}`);
    }
};

// const ROOT_URL = 'http://localhost:10001';

export async function getAllEndpoints(dispatch) {
    const requestOptions = {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' },
    };

    try {
        dispatch({ type: 'ENDPOINTS_REQUEST' });
        let response = await fetch(`/api/v1/mngt/`, requestOptions);
        let data = await response.json();
        if (data.length >= 0) {

            dispatch({ type: 'ENDPOINTS_SUCCESS', payload: data });
            return data;
        }

        if (data.hasOwnProperty('error')) {

            dispatch({ type: 'ENDPOINTS_ERROR', error: data });
            return;
        }
    } catch (error) {
        dispatch({ type: 'ENDPOINTS_ERROR', error: error });
        console.log(error);
    }
}

