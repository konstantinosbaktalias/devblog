import { useState } from 'react';
import { Button, Container, TextField } from '@material-ui/core'
import { Alert } from '@material-ui/lab'

import API from '../axios';

const CreatePost = () => {
    const [title, setTitle] = useState('');
    const [context, setContext] = useState('');
    const [error, setError] = useState(null);

    const createPost = async () => {
        await API.post('/create/post', {
            title: title,
            context: context
        }).then(res => {
            window.location.replace('/');
        }).catch(err => {
            setError(err.response.data);
        });
    }

    return (
        <Container style={{marginTop: "18px"}, {textAlign: "center"}}>
            <h1>Create post</h1>
            {(error != null) ?
                <Alert severity="error" style={{backgroundColor: "rgb(255, 0, 0, 0.2)", margin: "15px 23%"}}>{error}</Alert>
                :
                <></>
            }
            <form>
                <TextField label="Title" style={{width: "50%"}} multiline onChange={e =>setTitle(e.target.value)} /><br/><br/><br/>
                <TextField label="Context" style={{width: "50%"}} multiline rows={10} onChange={e => setContext(e.target.value)}/><br/><br/><br/>
                <Button onClick={createPost} variant="contained" style={{backgroundColor: "#38d100", border: "0"}}>Create</Button>
            </form>
        </Container>
    )
}

export default CreatePost;