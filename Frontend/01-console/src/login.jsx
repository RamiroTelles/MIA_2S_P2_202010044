import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
  } from "./components/ui/card"

import { Input } from "./components/ui/input";
import {Button} from './components/ui/button';
import { useState ,useEffect} from 'react';
import {useHistory} from "react-router-dom";

const Login = () =>{
    const backendUrl = import.meta.env.VITE_BACKEND_URL;

    const [user,setUser] = useState('');
    const [pass,setPass] = useState('');
    const [id,setId] = useState('');
    const [isPending,setIsPending] = useState(false);
    const [isLogin,setIsLogin] = useState(false);
    const history = useHistory()


   
    const postLogin = async () => {
        const data = {user,pass,id};
        setIsPending(true);

        fetch(backendUrl+'/login', {
            method: 'POST',
            headers: {"Content-type":"application/json"},
            body: JSON.stringify(data)
        }).then((res) => {
            if (!res.ok){
                console.log("Error al Loguearse",res);
                return
            }
            setIsPending(false)
            setIsLogin(true)
            alert("Logueo Correcto");
            
            history.push('/')
        });
    }

    const getLogout = async () =>{
        setIsPending(true);
        fetch(backendUrl+"/logout")
            .then((res) => {
                if(!res){
                    console.log("Error al desloguearse");
                    return
                }
                setIsPending(false);
                setIsLogin(false);
                alert("Deslogueado con exito");
                history.push("/")
            });
    }

    useEffect ( ()=>{
        fetch(backendUrl+'/getLogin')
            .then((res) => {
                //console.log(res);
                setIsLogin(res.ok);    
            });
    },[]);

    

    return(

        <>
            <Card >
                <CardHeader>
                    <CardTitle>Login</CardTitle>
                </CardHeader>
                <CardContent className="flex justify-center">
                {isLogin ? ( 
                <div>
                    <h4>Estas Logueado</h4>
                    {!isPending && <Button className='m-2 w-32 text-xl' onClick={getLogout}> Logout</Button>}
                    {isPending && <Button className='m-2 w-32 text-xl' disabled> Cargando...</Button>}  
                    <Button onClick={() => console.log(isLogin)}>asdasd</Button>
                </div>
                ) : (
                <div className=''>
                    <h4 className='flex'>Usuario </h4>
                    <Input value={user} onChange={(e) => setUser(e.target.value)}/>
                    
                    <h4 className='flex'>Contrasena</h4>
                    <Input type='password' value={pass} onChange={(e) => setPass(e.target.value)}/>
                    <h4 className='flex'>Id</h4>
                    <Input value={id} onChange={(e) => setId(e.target.value)}/>
                    <div >
                        {!isPending && <Button className='m-2 w-32 text-xl' onClick={postLogin}>Login</Button>}
                        {isPending && <Button className='m-2 w-32 text-xl' disabled>Cargando...</Button>}
                        <Button onClick={() => console.log(isLogin)}>asdasd</Button>
                    </div>
                </div>)}
                

                </CardContent>
                <CardFooter>
            
                </CardFooter>
            </Card>
        
        </>
    );

}


export default Login;