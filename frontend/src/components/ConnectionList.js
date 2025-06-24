import ConnectionControls from './ConnectionControls';
// ...existing code...
const [showHalfConnections, setShowHalfConnections] = useState(false);

useEffect(() => {
    fetch(`/api/connections?half=${showHalfConnections}`)
        .then(res => res.json())
        .then(setConnections);
}, [showHalfConnections]);

// Helper function to get connection class based on status
const getConnectionClass = (connection) => {
    if (!connection.status || connection.status === 'complete') {
        return 'connection-item';
    }
    return `connection-item half-connection ${connection.status}`;
};

// Helper component for status indicator
const StatusIndicator = ({ status }) => {
    if (!status || status === 'complete') {
        return null;
    }
    
    const label = status === 'request_only' ? 'Request Only' : 'Response Only';
    
    return (
        <div className={`status-indicator ${status}`}>
            {label}
        </div>
    );
};

return (
    <div className="connections-container">
        <ConnectionControls 
            showHalfConnections={showHalfConnections} 
            setShowHalfConnections={setShowHalfConnections} 
        />
        <div className="connections-list">
            {connections.map(connection => (
                <div 
                    key={connection.id}
                    className={getConnectionClass(connection)}
                    onClick={() => onSelectConnection(connection)}
                >
                    <div className="connection-header">
                        <span className="protocol">{connection.protocol}</span>
                        <span className="timestamp">{formatTimestamp(connection.timestamp)}</span>
                        <StatusIndicator status={connection.status} />
                    </div>
                    <div className="connection-details">
                        <div className="source-target">
                            {connection.source} → {connection.target}
                        </div>
                    </div>
                </div>
            ))}
        </div>
    </div>
);
