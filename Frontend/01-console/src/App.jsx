
import './App.css';
import {BrowserRouter as Router, Route,Switch} from 'react-router-dom';
import Console from './console'
import NavigBar from './NavigBar';
import Login from './login';
import Visualizer from './visualizer';


function App() {
  return (
    <>
    <Router>
      <div className='App'>
        <NavigBar />
      </div>
        <div className='bodyContent'>    
            

              <Switch>
                <Route exact path='/'>
                  <Console />
                </Route>
                <Route path='/login'>
                  <Login />
                </Route>
                <Route path='/visualizer'>
                  <Visualizer />
                </Route>
              </Switch>
            
          
       </div>
       </Router>
    </>
  )
}

export default App;
