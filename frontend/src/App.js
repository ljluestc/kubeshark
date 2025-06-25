function App() {
  const [connections, setConnections] = useState([]);
  const [showHalfConnections, setShowHalfConnections] = useState(
    process.env.REACT_APP_SHOW_HALF_CONNECTIONS === "true"
  );
  
  // Existing code...

  return (
    <div className="app">
      <Header />
      <main>
        <ConnectionControls 
          showHalfConnections={showHalfConnections} 
          setShowHalfConnections={setShowHalfConnections} 
        />
        <ConnectionList 
          connections={connections} 
          showHalfConnections={showHalfConnections} 
        />
      </main>
    </div>
  );
}
