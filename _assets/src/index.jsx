import 'babel-polyfill'

import React from 'react'
import { render } from 'react-dom'
import { Provider } from 'react-redux'
import { createStore, applyMiddleware } from 'redux'
import thunkMiddleware from 'redux-thunk'

import reducers from './app/reducers'
import App from './app/App'

let store = createStore(
    reducers,
    applyMiddleware(
        thunkMiddleware
    )
)

render(
    <Provider store={store}>
        <App />
    </Provider>,
    document.getElementById('root')
)