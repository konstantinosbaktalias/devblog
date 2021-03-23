import API from '../axios';

const DeletePost = async(id, path) => {
    await API.post(`/delete/post/${id}`).then(res => {
        window.location.replace(path);        
    }).catch(err => {
        console.log(err);
    })
}

export default DeletePost;