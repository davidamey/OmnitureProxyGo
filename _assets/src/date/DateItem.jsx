import React, { PropTypes } from 'react'

class DateItem extends React.Component {
    static propTypes = {
        onClick: PropTypes.func.isRequired,
        date: PropTypes.string.isRequired,
        selected: PropTypes.bool
    }

    state = {
        hover: false
    }

    render() {
        return (
            <li
                onMouseEnter={() => { this.setState({ hover: true }) } }
                onMouseLeave={() => { this.setState({ hover: false }) } }
                onClick={this.props.onClick}
                className="date"
                style={{
                    background: (this.props.selected || this.state.hover) ? 'lightgreen' : 'none',
                }}
            >
                { this.props.date }
            </li>
        );
    }
};

// Date.propTypes = {
//     onClick: PropTypes.func.isRequired,
//     date: PropTypes.string.isRequired,
//     selected: PropTypes.bool
// }

export default DateItem