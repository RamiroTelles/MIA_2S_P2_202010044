import {
    Card,
    CardContent,

  } from "./components/ui/card";

import DriveImage from "./imgs/Drive-disk.png"
import FileImage from "./imgs/File.png"
import FolderImage from "./imgs/Folder.png"
import PartitionImage from "./imgs/Partition.png"


const CardFiles = (props) => {
    const files = props.files;
    const changePath = props.changePath

    return (
        <>
            {files.map((file) => {
                // Selección de la imagen según el valor de `file.ls_type`
                let imgSrc;
                switch (file.ls_type) {
                    case 0:
                        imgSrc = DriveImage;
                        break;
                    case 1:
                        imgSrc = PartitionImage;
                        break;
                    case 2:
                        imgSrc = FolderImage;
                        break;
                    case 3:
                        imgSrc = FileImage;
                        break;
                    default:
                        imgSrc = FileImage; // Imagen predeterminada si no coincide
                        break;
                }

                return (
                    <Card key={file.ls_name} id={file.ls_type} onClick={() => changePath(file.ls_name,file.ls_type)}>
                        <CardContent className="flex flex-col items-center">
                            <img src={imgSrc} alt={file.ls_name} className="w-auto h-auto" />
                            <div className="text-center">{file.ls_name}</div>
                        </CardContent>
                    </Card>
                );
            })}
        </>
    );
};

export default CardFiles