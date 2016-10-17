// Component
import React, { PropTypes } from 'react'
import Visitor from './Visitor.jsx'

const VisitorList = ({ visitors, onVisitorClick }) => (
    <ul>
        { visitors.map( vid =>
            <Visitor
                key={vid}
                vid={vid}
                onClick={ () => onVisitorClick(vid) }
            />
        )}
    </ul>
);

VisitorList.propTypes = {
    visitors : PropTypes.arrayOf(PropTypes.string.isRequired).isRequired,
    onVisitorClick: PropTypes.func.isRequired
};

export { VisitorList }

// Container
import { connect } from 'react-redux'
import { selectVisitor } from './VisitorActions'

const mapStateToProps = (state) => {
    return {
        visitors : Object.keys(state.logs[state.selectedDate] || {})
    }
}

const mapDispatchToProps = (dispatch) => {
    return {
        onVisitorClick: (vid) => {
            // alert(`visitor ${id} clicked`);
            dispatch(selectVisitor(vid))
        }
    }
}

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(VisitorList)