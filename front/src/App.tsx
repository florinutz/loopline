import React, {Component} from 'react';
import {Note} from "./models/note";
import {NewNoteForm} from "./components/NewNoteForm";
import {NoteList} from "./components/NoteList";

interface State {
    newNote: Note;
    notes: Note[];
}

const uuid = require('uuid/v4');

class App extends Component<{}, State> {
    state: State = {
        newNote: {id: uuid(), title: "title", content: "content", created_date: new Date()},
        notes: []
    };

    componentDidMount() {
        let parent = this;

        let callback = function () {
            fetch("http://localhost:8080/", {method: 'GET'}).then(response => response.json()).then(
                (result) => {
                    if (!result) {
                        return
                    }
                    parent.setState(state => {
                        return ({...state, notes: [...result]});
                    });
                },
                (error) => console.log("can't load notes from backend: ", error)
            )
        };

        callback();

        // setInterval(callback, 2000);
    }

    render = () => (
        <div>
            <NewNoteForm
                onChange={this.onNewTitleChange}
                onContentChange={this.onNewContentChange}
                onAdd={this.onAdd}
                note={this.state.newNote}
            />
            <NoteList notes={this.state.notes} onDelete={this.onDelete}/>
        </div>
    );

    private onAdd = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        this.setState(state => ({
            notes: [state.newNote, ...state.notes] // add new note above the others
        }));

        fetch("http://localhost:8080/", {
                method: 'POST',
                body: JSON.stringify(this.state.newNote)
            }
        ).catch(error => console.log("can't delete note:", error));
    };

    private onNewTitleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        this.setState({
            newNote: {
                ...this.state.newNote,
                title: event.target.value,
            }
        });
    };

    private onNewContentChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
        this.setState({
            newNote: {
                ...this.state.newNote,
                content: event.target.value,
            }
        });
    };

    private onDelete = (noteToDelete: Note) => {
        this.setState(previousState => ({
            notes: [...previousState.notes.filter(note => note.id !== noteToDelete.id)]
        }));

        fetch("http://localhost:8080/", {
                method: 'DELETE',
                body: JSON.stringify({"ids": [noteToDelete.id]})
            }
        ).catch(error => console.log("can't delete note:", error));
    };
}

export default App;
