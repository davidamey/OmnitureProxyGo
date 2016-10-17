// var React = require('react');
// var ReactDOM = require('react-dom');

import React from 'react'
import ReactDOM from 'react-dom'

import VisitorList from './visitor/VisitorList.jsx'

// import { createStore } from 'redux'

// const VisitorComponent = vid => ( <div className="visitor">vid={ vid }</div> );

// const VisitorList = vids => (
//     <div>
//         vids= { vids.map( vid => VisitorComponent(vid) ) }
//     </div>
// );

// ReactDOM.render(VisitorComponent(1), document.getElementById('visitor'));
// ReactDOM.render(VisitorList([1, 2, 3]), document.getElementById('visitors'));
ReactDOM.render(
    <VisitorList
        vids={ ['1', '2', '3'] }
        onVisitorClick={ (vid) => console.log(vid) }
    />, document.getElementById('visitors'));

// let store = createStore()