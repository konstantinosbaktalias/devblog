import { useState } from 'react';
import { Button, Container, TextField } from '@material-ui/core'
import { Alert } from '@material-ui/lab'

import API from '../axios';

const Login = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState(null);

    const LoginUser = async () => {
        await API.post('/user/login', {
            username: username,
            password: password
        }).then(res => {
            window.location.replace('/');            
        }).catch(err => {
            setError(err.response.data);
        });
    }

    return (
        <Container style={{marginTop: "18px"}, {textAlign: "center"}}>
            <h1>Login</h1>
            {(error != null) ?
                <Alert severity="error" style={{backgroundColor: "rgb(255, 0, 0, 0.2)", margin: "15px 23%"}}>{error}</Alert>
                :
                <></>
            }
            <form>
                <TextField label="Username" style={{width: "50%"}} onChange={e =>setUsername(e.target.value)} /><br/><br/><br/>
                <TextField type="password" label="Password" style={{width: "50%"}} onChange={e => setPassword(e.target.value)}/><br/><br/><br/>
                <Button onClick={LoginUser} variant="contained" color="primary">Login</Button>
            </form>
        </Container>
    )
}

export default Login;