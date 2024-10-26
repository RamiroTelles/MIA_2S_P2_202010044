import {
    NavigationMenu,
    NavigationMenuContent,
    NavigationMenuIndicator,
    NavigationMenuItem,
    NavigationMenuLink,
    NavigationMenuList,
    NavigationMenuTrigger,
    
    navigationMenuTriggerStyle,
    NavigationMenuViewport,
  } from "./components/ui/navigation-menu"

import {Link} from 'react-router-dom'


const NavigBar = () =>{
    return(
        <NavigationMenu className="pb-2">
            <NavigationMenuList>
                <NavigationMenuItem>
                    <Link to="/" className={navigationMenuTriggerStyle()} >
                        
                            Consola
                        
                    </Link>
                </NavigationMenuItem>
                <NavigationMenuItem>
                    <Link to="/login" className={navigationMenuTriggerStyle()} >
                        
                            Login
                       
                    </Link>
                </NavigationMenuItem>
                <NavigationMenuItem>
                    <Link to="/visualizer" className={navigationMenuTriggerStyle()} >
                        
                            Visualizador
                       
                    </Link>
                </NavigationMenuItem>
            </NavigationMenuList>
        </NavigationMenu>
    );
}


export default NavigBar;