import React, {FunctionComponent} from "react";
import {Note} from "../models/note";

interface Props {
    note: Note,
    onDelete: (note: Note) => void;
}

export const NoteListItem: FunctionComponent<Props> = ({note, onDelete}) => {
    const onClick = () => {
        onDelete(note)
    };

    return (
        <li>
            <h4>{note.title} <button onClick={onClick}>X</button></h4>
            <p>{note.content}</p>
        </li>
    )
};
