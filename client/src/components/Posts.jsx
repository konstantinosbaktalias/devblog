import { useState, useEffect } from 'react';
import { Link, useParams } from 'react-router-dom';
import { Container, Paper, Button, ButtonGroup, LinearProgress } from '@material-ui/core'
import { Pagination, PaginationItem } from '@material-ui/lab';

import API from '../axios';

import DeletePost from './DeletePostFunc';

const Posts = ({loggedUser}) => {
    let { page } = useParams();
    page = parseInt(page);

    const [posts, setPost] = useState([]);
    const [totalPages, setTotalPages] = useState(0);

    const [loaded, setLoaded] = useState(false);

    const getPosts = async(pg) => {
        await API.get(`/posts/pages/${pg}`).then(res => {
            setPost(res.data.posts);
            setTotalPages(res.data.totalPages);
        }).catch(err => {
            console.log(err)
        });
        setLoaded(true);
    };

    useEffect(() => {
        getPosts(page)
    }, [page]);
    
    return (
        <>
        {(!loaded)?
            <LinearProgress color="secondary" />
            :
            <Container style={{marginTop: "18px"}}>   
                {(loggedUser != null) ?
                    <>
                        <Link style={{ textDecoration: "none"}} to={`/create/post`}><Button variant="contained" style={{backgroundColor: "#38d100", border: "0", margin: "12px", float: "right"}}>Create</Button></Link>
                        <br/><br/><br/><br/>
                    </>
                    :
                    <></>
                }
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
                                        <Button style={{margin: "0 15px"}} variant="contained" onClick={() => DeletePost(post.Id, `/posts/${page}`)} color="secondary">Delete</Button>
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
                        to={`/posts/${item.page}`}
                        {...item}/>
                    )}/>
                </div>
            </Container>
        }
        </>
    )
}

export default Posts;