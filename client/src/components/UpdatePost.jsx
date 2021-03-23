import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { Button, Container, TextField } from '@material-ui/core'
import { Alert } from '@material-ui/lab'

import API from '../axios';

const UpdatePost = () => {
    let { id } = useParams();

    const [title, setTitle] = useState('');
    const [context, setContext] = useState('');
    const [error, setError] = useState(null);

    const [post, setPost] = useState([]);

    const getPost = async(id) => {
        await API.get(`/post/${id}`).then(res => {
            setPost(res.data);
        }).catch(err => {
            setError(err.response)
            console.log(err)
        });
    };

    const updatePost = async(id) => {
        await API.post(`/update/post/${id}`, {
            title: title,
            context: context
        }).then(res => {
            window.location.replace('/');
        }).catch(err => {
            setError(err.response.data);
        });
    }

    useEffect(() => {
        getPost(id)
    }, [id])


    return (
        <Container style={{marginTop: "18px"}, {textAlign: "center"}}>
            {(error != null) ?
                <Alert severity="error" style={{backgroundColor: "rgb(255, 0, 0, 0.2)", margin: "15px 23%"}}>{error}</Alert>
                :
                <></>
            }
            <form>
                <h1>Update post</h1>
                <TextField label="Title" style={{width: "50%"}} multiline defaultValue={post.Title} onChange={e =>setTitle(e.target.value)} /><br/><br/><br/>
                <TextField label="Context" style={{width: "50%"}} multiline rows={10} defaultValue={post.Context} onChange={e => setContext(e.target.value)}/><br/><br/><br/>
                <Button onClick={() => updatePost(id)} variant="contained" style={{backgroundColor: "#dec800", border: "0"}}>Update post</Button>
            </form>
        </Container>
    )
}

export default UpdatePost;