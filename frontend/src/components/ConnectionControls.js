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
    </div>
  );
};

export default ConnectionControls;
