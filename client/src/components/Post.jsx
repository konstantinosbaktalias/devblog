import { useState, useEffect } from 'react';
import { Link, useParams } from 'react-router-dom';
import { Container, Paper, Button, ButtonGroup, LinearProgress } from '@material-ui/core'

import API from '../axios';

import DeletePost from './DeletePostFunc';

const Post = ({loggedUser}) => {
    let { id } = useParams();

    const [post, setPost] = useState([]);
    const [error, setError] = useState(null);

    const [loaded, setLoaded] = useState(false);

    const getPost = async(id) => {
        await API.get(`/post/${id}`).then(res => {
            setPost(res.data);
        }).catch(err => {
            setError(err.response)
            console.log(err)
        });
        setLoaded(true);
    };

    useEffect(() => {
        getPost(id)
    }, [id])
    
    return (
        <>
        {(!loaded)?
            <LinearProgress color="secondary" />
            :
            <Container style={{marginTop: "18px"}}>   
            {(error != null) ?
                    <h1>{error.status} {error.data}</h1>
                :
                <Paper elevation={5} p={2} square>
                    <div style={{padding: "20px"}}>
                        <h2>{post.Title}</h2>
                        <Link style={{ textDecoration: "none", color: "#fff"}} to={`/user/${post.Author}`}><span style={{color: "#555"}}>@{post.Author}</span></Link>
                        <p>{post.Context}</p>
                        {(loggedUser != null && loggedUser.Username == post.Author) ?
                            <ButtonGroup>
                                <Link style={{ textDecoration: "none"}} to={`/update/post/${post.Id}`}><Button style={{backgroundColor: "#dec800", border: "0"}}>Update</Button></Link>
                                <Button style={{margin: "0 15px"}} onClick={() => DeletePost(post.Id, '/')} variant="contained" color="secondary">Delete</Button>
                            </ButtonGroup>
                            :
                            <></>
                        }
                    </div>
                </Paper>
            }
            </Container>
        }   
        </>
    )
}

export default Post;