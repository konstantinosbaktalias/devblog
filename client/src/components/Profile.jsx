import { useState, useEffect } from 'react';
import { Link, useParams } from 'react-router-dom';
import { Container, Paper, Button, ButtonGroup, LinearProgress } from '@material-ui/core'
import { Pagination, PaginationItem } from '@material-ui/lab';

import API from '../axios';

import DeletePost from './DeletePostFunc';

const Profile = ({loggedUser}) => {
    let { username, page } = useParams();
    page = parseInt(page)

    const [posts, setPost] = useState([]);
    const [totalPages, setTotalPages] = useState(0);

    const [loaded, setLoaded] = useState(false);
    const [error, setError] = useState(null);

    const getUserProfile = async(username, page) => {
        await API.get(`/users/${username}/${page}`).then(res => {
            setPost(res.data.posts);
            setTotalPages(res.data.totalPages);
            console.log(res)
        }).catch(err => {
            setError(err.response)
            console.log(err.response)
        });
        setLoaded(true);
    };

    useEffect(() => {
        getUserProfile(username, page);
    }, [username, page]);
    
    return (
        <>
        {(!loaded)?
            <LinearProgress color="secondary" />
            :
            <Container style={{marginTop: "18px"}}>   
                {(error != null) ?
                    <h1>{error.status} {error.data}</h1>
                    :
                    <>
                    <h2>{username}'s posts:</h2>
                    {(posts == null) ? <h4>Posts not avialable</h4> : posts.map(post => {
                        return (
                            <div style={{margin: "5px"}}>
                                <Paper elevation={5} p={2} square>
                                    <div style={{padding: "20px"}}>
                                    <Link style={{ textDecoration: "none", color: "#000"}} to={`/post/${post.Id}`}><h3>{post.Title}</h3></Link>
                                    <Link style={{ textDecoration: "none", color: "#fff"}} to={`/user/${post.Author}/1`}><span style={{color: "#555"}}>@{post.Author}</span></Link>
                                        <p>{post.Context.substring(0, 100)}{(post.Context.length > 100) ? <>...</> : <></>}</p>
                                        {(loggedUser != null && loggedUser.Username == post.Author) ?
                                        <ButtonGroup>
                                            <Link style={{ textDecoration: "none"}} to={`/update/post/${post.Id}`}><Button style={{backgroundColor: "#dec800", border: "0"}}>Update</Button></Link>
                                            <Button style={{margin: "0 15px"}} variant="contained" onClick={() => DeletePost(post.Id, `/user/${username}/${page}`)} color="secondary">Delete</Button>
                                        </ButtonGroup>
                                        :
                                        <></>
                                        }
                                    </div>
                                </Paper>
                            </div>
                        )
                    })}
                    <div style={{margin: "25px 0"}}>
                    <Pagination 
                    count={totalPages} 
                    color="primary" 
                    page={page} 
                    shape="rounded" 
                    showLastButton
                    renderItem={item => (
                        <PaginationItem
                        component={Link}
                        to={`/user/${username}/${item.page}`}
                        {...item}/>
                    )}/>
                    </div>
                    </>
                }
            </Container>
        }
        </>
    )
}

export default Profile;