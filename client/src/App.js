import { useState, useEffect } from 'react';
import { BrowserRouter as Router, Switch, Route, Redirect, Link } from 'react-router-dom'; 
import { AppBar, Button, Toolbar, Typography } from '@material-ui/core'

import API from './axios';

import CreatePost from './components/CreatePost.jsx';
import UpdatePost from './components/UpdatePost.jsx';
import Posts from './components/Posts.jsx';
import Post from './components/Post.jsx';
import Login from './components/Login.jsx';
import Signup from './components/Signup.jsx';
import Profile from './components/Profile.jsx';

const App = () => {
  const [state1, setState] = useState(1);
  const [loggedUser, setLoggedUser] = useState(null);

  const getUser = async () => {
    await API.get('/user/me').then(res => {
      setLoggedUser(res.data)
    }).catch(err => {
      console.log(err.response)
    })
    console.log(loggedUser)
  }

  const LogoutUser = () => {    
    document.cookie = 'auth_token=; path=/; Max-Age=0';
    setLoggedUser(null);
    window.location.replace('/');
  }

  useEffect(() => {
    getUser()
  }, []);  

  return (
    <Router>
      <AppBar position="static">
        <Toolbar>
          <Link style={{ textDecoration: "none", color: "#fff"}} to="/"><Typography variant="h6">DevBlog</Typography></Link>
          {(loggedUser == null) ?
            <div style={{width: "100%"}}>
              <Link style={{ textDecoration: "none", color: "#fff"}} to="/login"><Button style={{paddingLeft: "15px"}} style={{float: "right"}}  color='inherit'>Login</Button></Link>
              <Link style={{ textDecoration: "none", color: "#fff"}} to="/signup"><Button style={{paddingLeft: "15px"}} style={{float: "right"}} color='inherit'>Signup</Button></Link>
            </div>
            :
            <div style={{width: "100%"}}>
              <Button  onClick={LogoutUser} style={{paddingLeft: "15px", color: "#fff"}} style={{float: "right"}}  color='inherit'>Logout</Button>
              <Link style={{ textDecoration: "none", color: "#fff"}} to={`/user/${loggedUser.Username}/1`}><Button style={{paddingLeft: "15px"}} style={{float: "right"}} color='inherit'>Profile</Button></Link>
            </div>
          }
        </Toolbar>
      </AppBar>

      <Switch>
        <Route path='/posts/:page/' >
          <Posts loggedUser={loggedUser}/>
        </Route>
        <Route path="/post/:id">
          <Post loggedUser={loggedUser}/>
        </Route>
        <Route path="/create/post">
          <CreatePost/>
        </Route>
        <Route path="/update/post/:id">
          <UpdatePost />
        </Route>
        <Route path='/login'>
          <Login />
        </Route>
        <Route path='/signup'>
          <Signup/>
        </Route>
        <Route path='/user/:username/:page'>
          <Profile loggedUser={loggedUser}/>
        </Route>
        <Route path='/'>
          <Redirect to='/posts/1'/>
        </Route>
      </Switch>
    </Router>
  );
}

export default App;