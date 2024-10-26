import {
    Card,
    CardContent,

    CardHeader,
    CardTitle,
  } from "./components/ui/card";

import { useEffect, useState } from 'react';

import { Input } from "./components/ui/input";
import {Button} from './components/ui/button';
import { AspectRatio } from "./components/ui/aspect-ratio"

import CardFiles from "./cardFiles";

const Visualizer = () =>{
    const backendUrl = import.meta.env.VITE_BACKEND_URL;
    const [Path, setPath] = useState('/');
    
    const [files,setFiles] = useState(null);
    const [pathDisco,setPathDisco] = useState('');
    const [nombreParticion,setNombreParticion] = useState('');
    //const [flagArchivo,setFlagArchivo] = useState(false);
    const [prevState,setPrevState] = useState([]);

  
   
    const changePath = (name,type)=>{
        
        setPrevState(prevState => [...prevState, { pathDisco, nombreParticion, Path }]);
        
        if (type ==0){
            setPathDisco(name);
        }else if (type==1){
            setNombreParticion(name);
        }else if(type==2){
            setPath(Path+name+"/");
        }
    }

    const goBack = () =>{
        if (prevState.length>0){
            setPathDisco(prevState[prevState.length-1].pathDisco);
            setNombreParticion(prevState[prevState.length-1].nombreParticion);
            setPath(prevState[prevState.length-1].Path);

            setPrevState(prevState => prevState.slice(0, -1)); 
        }else{
            console.log("No se puede volver atras");
        }
        
        
    }

    const ver = () =>{
        console.log(pathDisco);
        console.log(nombreParticion);
        console.log(Path);
        console.log(prevState);
    }

    useEffect( ()=>{
        const data = {pathDisco,nombreParticion,Path}

        fetch(backendUrl+"/ls",{
            method: 'POST',
            headers: {"Content-type":"application/json"},
            body: JSON.stringify(data)
        }).then((res) => {
            if(!res.ok){
                throw Error("No se obtuvo la informacion");
            }
            return res.json();
           
        }).then((data) => {
            console.log(data);
            setFiles(data)
        })

    },[pathDisco,nombreParticion,Path])

    return(
    <>
        <Card>
            <CardHeader>
                <CardTitle>Visualizador</CardTitle>
            </CardHeader>
            <CardContent>
                <div className="p-4">
                    <Button onClick={ver}>Ver</Button>
                    <Button className="flex mb-2" onClick={goBack}>Atras</Button>
                    <Input className="mb-2" value={Path} onChange={(e) => setPath(e.target.value) }/>
                    <AspectRatio ratio={16/9} className="grid grid-cols-7 grid-rows-4 gap-1 border bg-muted">
                    {files && <CardFiles files={files} changePath={changePath}></CardFiles>}
                        
                    </AspectRatio>  
                </div>
                
            </CardContent>
           


        </Card>
    
    
    
    </>
    
);

}

export default Visualizer;