import React, {FunctionComponent} from "react";
import {Note} from "../models/note";
import {NoteListItem} from "./NoteListItem";

interface Props {
    notes: Note[];
    onDelete: (note: Note) => void;
}

export const NoteList: FunctionComponent<Props> = ({notes, onDelete}) => {
    return (notes.length > 0) ?
        <ul>{notes.map(note =>
            <NoteListItem note={note} onDelete={onDelete} key={note.id}/>
        )}</ul>
        :
        <p>no items</p>;
};
