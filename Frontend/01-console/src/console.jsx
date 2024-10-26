import { useState } from 'react'

import './App.css'
import {Button} from './components/ui/button'
import {Textarea} from './components/ui/textarea'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "./components/ui/card"


function Console() {
  const backendUrl = import.meta.env.VITE_BACKEND_URL;
  const [consoleText, setConsoleText] = useState(''); // Estado para almacenar el texto del Textarea
  const [inputText, setInputText] = useState('');

  const reiniciarAPI = ()=>{
    fetch(backendUrl+'/restartValues')
    .then((res) => {
      if (res.ok){
        alert("Reiniciado con exito");
      }
    })
  }
  const handleInputText = (event) => {

    setInputText((event.target.value));
  }

  // Función que maneja el clic del botón para agregar texto
  const handleAddText = async () => {
    // Actualiza el texto del Textarea agregando nuevo contenido
    
    try{
      
      const response = await fetch(backendUrl+ '/analizar', {
        method: 'POST',
        headers: {'Content-Type':'application/json'},
        body: JSON.stringify({
          Cmd: inputText
        })
        
      });
      
      if(response.ok){
        const jsonResponse = await response.json();
        
        const {result} = jsonResponse;
        //console.log(result)

        setConsoleText((consoleText) => result);
      }else{
        console.log(response.status)
      }

    }catch(error){
      console.log(error)
    }
    
    //setConsoleText((consoleText) => consoleText + inputText); // Agrega el texto deseado
  };

  return (
    

    <>
      
        
        <Card>
          <CardHeader>
            <CardTitle>Proyecto1</CardTitle>
            
          </CardHeader>
          <CardContent className="grid gap-4">
            <div className=''>
              <h4 className="flex">Entrada</h4>
              <Textarea className='m-1 h-64' value={inputText} onChange={handleInputText}/>
              <div className='flex'>
                <Button className='m-2' onClick={handleAddText}>Ejecutar</Button>
                <Button className='m-2' onClick={reiniciarAPI}>Reinciar</Button>
              </div>
              <h4 className='flex'>Consola</h4>
              <Textarea className='m-1 h-64' value={consoleText} readOnly/>
            </div>
          </CardContent>
          <CardFooter>
            
          </CardFooter>
        </Card>

    </>
  )
}

export default Console;