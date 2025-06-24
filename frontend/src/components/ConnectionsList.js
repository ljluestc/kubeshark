// ...existing code...
const [showHalfConnections, setShowHalfConnections] = useState(false);

useEffect(() => {
    fetch(`/api/connections?half=${showHalfConnections}`)
        .then(res => res.json())
        .then(setConnections);
}, [showHalfConnections]);

// Get connection status indicator component
const ConnectionStatusIndicator = ({ status }) => {
    if (status === 'complete') {
        return null; // No indicator needed for complete connections
    }
    
    return (
        <div className={`connection-status ${status}`}>
            {status === 'request_only' ? 'Request Only' : 'Response Only'}
        </div>
    );
};

return (
    <div className="connections-container">
        <div className="connections-controls">
            <label className="half-connections-toggle">
                <input
                    type="checkbox"
                    checked={showHalfConnections}
                    onChange={(e) => setShowHalfConnections(e.target.checked)}
                />
                Show Half-Connections
            </label>
        </div>
        <div className="connections-list">
            {connections.map(connection => (
                <div 
                    key={connection.id}
                    className={`connection-item ${connection.status !== 'complete' ? 'half-connection' : ''}`}
                    onClick={() => onSelectConnection(connection)}
                >
                    <div className="connection-header">
                        <span className="protocol">{connection.protocol}</span>
                        <span className="timestamp">{formatTimestamp(connection.timestamp)}</span>
                        <ConnectionStatusIndicator status={connection.status} />
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
