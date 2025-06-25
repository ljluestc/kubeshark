import React from 'react';
import './ConnectionItem.css';

class ConnectionItem extends React.Component {
    getConnectionClass = (connection) => {
        if (!connection.status || connection.status === 'complete') {
            return 'connection-item';
        }
        return `connection-item half-connection ${connection.status}`;
    }

    render() {
        const { connection } = this.props;
        const connectionClass = this.getConnectionClass(connection);
        
        return (
            <div className={connectionClass}>
                <div className="connection-header">
                    <span className="connection-id">{connection.id}</span>
                    <span className="connection-timestamp">{new Date(connection.timestamp).toLocaleTimeString()}</span>
                    {connection.status && connection.status !== 'complete' && (
                        <span className="connection-status">{connection.status}</span>
                    )}
                </div>
                <div className="connection-details">
                    <div className="connection-source">{connection.source}</div>
                    <div className="connection-arrow">→</div>
                    <div className="connection-target">{connection.target}</div>
                </div>
                <div className="connection-protocol">{connection.protocol}</div>
            </div>
        );
    }
}

export default ConnectionItem;
