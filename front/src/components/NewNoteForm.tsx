import React, {FunctionComponent} from "react";
import {Note} from "../models/note";

interface Props {
    onChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
    onContentChange: (event: React.ChangeEvent<HTMLTextAreaElement>) => void;
    onAdd: (event: React.FormEvent<HTMLFormElement>) => void;
    note: Note;
}

export const NewNoteForm: FunctionComponent<Props> = ({onChange, onContentChange, onAdd, note}) => (
    <form onSubmit={onAdd}>
        <input onChange={onChange} placeholder={note.title} /><br/>
        <textarea onChange={onContentChange} placeholder={note.content} /><br/>
        <button type="submit">Add a note</button>
    </form>
);
