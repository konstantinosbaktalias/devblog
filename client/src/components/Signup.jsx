import { useState } from 'react';
import { Button, Container, TextField } from '@material-ui/core'
import { Alert } from '@material-ui/lab'

import API from '../axios';

const Signup = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [confpassword, setConfPassword] = useState('');
    const [error, setError] = useState(null);

    const SignupUser = async () => {
        if(confpassword == password)
        {
            await API.post('/user/signup', {
                username: username,
                password: password
            }).then(res => {
                window.location.replace('/login');
            }).catch(err => {
                setError(err.response.data);
            });
        }
        else 
        {
            setError('Passwords do not match');
        }
    }

    return (
        <Container style={{marginTop: "18px"}, {textAlign: "center"}}>
            <h1>Signup</h1>
            {(error != null) ?
                <Alert severity="error" style={{backgroundColor: "rgb(255, 0, 0, 0.2)", margin: "15px 23%"}}>{error}</Alert>
                :
                <></>
            }
            <form>
                <TextField label="Username" style={{width: "50%"}} onChange={e =>setUsername(e.target.value)} /><br/><br/><br/>
                <TextField type="password" label="Password" style={{width: "50%"}} onChange={e => setPassword(e.target.value)}/><br/><br/><br/>
                <TextField type="password" label="Confirm password" style={{width: "50%"}} onChange={e => setConfPassword(e.target.value)}/><br/><br/><br/>
                <Button onClick={SignupUser} variant="contained" color="primary">Signup</Button>
            </form>
        </Container>
    )
}

export default Signup;