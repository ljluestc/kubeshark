import React from 'react';

const ConnectionControls = ({ showHalfConnections, setShowHalfConnections }) => {
  return (
    <div className="connection-controls">
      <label className="toggle-control half-connections-toggle">
        <input
          type="checkbox"
          checked={showHalfConnections}
          onChange={(e) => setShowHalfConnections(e.target.checked)}
        />
        <span className="toggle-label">Show Half-Connections</span>
      </label>
      <div className="control-item">
        <label className="checkbox-label">
          <input 
            type="checkbox" 
            checked={showHalfConnections} 
            onChange={e => setShowHalfConnections(e.target.checked)}
          />
          Show Half-Connections
        </label>
        <span className="tooltip">
          Display incomplete transactions where either the request or response is missing
        </span>
      </div>
    </div>
  );
};

export default ConnectionControls;
