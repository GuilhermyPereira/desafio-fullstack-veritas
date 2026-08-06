import React, { useState, useEffect } from 'react';
import { DragDropContext, Droppable, Draggable } from '@hello-pangea/dnd';
import axios from 'axios';
import { Edit2, Trash2, Plus, Moon, Sun } from 'lucide-react';
import './App.css';

const API_URL = 'http://localhost:8080/tasks';
const COLUMNS = ['A Fazer', 'Em Progresso', 'Concluídas'];

function App() {
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [formData, setFormData] = useState(null);

  // Estado para o Modo Escuro
  const [isDarkMode, setIsDarkMode] = useState(false);

  useEffect(() => {
    fetchTasks();
  }, []);

  const fetchTasks = async () => {
    try {
      setLoading(true);
      const response = await axios.get(API_URL);
      setTasks(response.data || []);
      setError(null);
    } catch (err) {
      setError('Erro ao carregar tarefas. Verifique se o backend está rodando.');
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async (e) => {
    e.preventDefault();
    try {
      if (formData.id) {
        await axios.put(`${API_URL}/${formData.id}`, formData);
      } else {
        await axios.post(API_URL, formData);
      }
      setFormData(null);
      fetchTasks();
    } catch (err) {
      setError('Erro ao salvar tarefa.');
    }
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Deseja realmente excluir esta tarefa?')) return;
    try {
      await axios.delete(`${API_URL}/${id}`);
      fetchTasks();
    } catch (err) {
      setError('Erro ao deletar tarefa.');
    }
  };

  const onDragEnd = async (result) => {
    const { destination, source, draggableId } = result;

    // Se soltou fora de uma coluna válida
    if (!destination) return;

    // Se soltou exatamente no mesmo lugar
    if (destination.droppableId === source.droppableId && destination.index === source.index) return;

    // Encontra a tarefa que foi arrastada
    const task = tasks.find(t => t.id === draggableId);
    const newStatus = destination.droppableId;

    // 1. Pega todas as tarefas que NÃO são da coluna de destino e não são a tarefa movida
    const otherTasks = tasks.filter(t => t.status !== newStatus && t.id !== draggableId);

    // 2. Pega as tarefas da coluna de destino na ordem atual (sem a tarefa movida)
    const destTasks = tasks.filter(t => t.status === newStatus && t.id !== draggableId);

    // 3. Insere a tarefa arrastada exatamente na nova posição (index) da coluna
    const updatedTask = { ...task, status: newStatus };
    destTasks.splice(destination.index, 0, updatedTask);

    // 4. Junta tudo para formar o novo estado mantendo a ordem visual correta
    setTasks([...otherTasks, ...destTasks]);

    // Envia a atualização para a API
    try {
      await axios.put(`${API_URL}/${draggableId}`, updatedTask);
    } catch (err) {
      setError('Erro ao mover tarefa. A alteração foi desfeita.');
      fetchTasks(); // Reverte em caso de erro
    }
  };


  if (loading) return <div className="loading">Carregando tarefas...</div>;

  return (
    <div className="app-container" data-theme={isDarkMode ? 'dark' : 'light'}>

      {/* Header com o Botão de Tema */}
      <header className="header" allign="center">
        <h1>Kanban Veritas</h1>
        <button
          className="theme-toggle"
          onClick={() => setIsDarkMode(!isDarkMode)}
        >
          {isDarkMode ? <Sun size={18} /> : <Moon size={18} />}
          {isDarkMode ? 'Modo Claro' : 'Modo Escuro'}
        </button>
      </header>

      {error && <div className="error">{error}</div>}

      {/* Modal de Formulário */}
      {formData && (
        <div className="modal-overlay">
          <form className="form-container" onSubmit={handleSave}>
            <h2>{formData.id ? 'Editar Tarefa' : 'Nova Tarefa'}</h2>
            <input
              required
              placeholder="Título da tarefa (obrigatório)"
              value={formData.title}
              onChange={e => setFormData({ ...formData, title: e.target.value })}
            />
            <textarea
              placeholder="Descrição (opcional)"
              value={formData.description}
              onChange={e => setFormData({ ...formData, description: e.target.value })}
            />
            <div className="form-actions">
              <button type="button" onClick={() => setFormData(null)}>Cancelar</button>
              <button type="submit">Salvar</button>
            </div>
          </form>
        </div>
      )}

      {/* Board e Colunas */}
      <DragDropContext onDragEnd={onDragEnd}>
        <div className="board">
          {COLUMNS.map(columnStatus => {
            const columnTasks = tasks.filter(t => t.status === columnStatus);

            return (
              <div className="column" key={columnStatus}>

                {/* Título fora da área de arrastar, com contador */}
                <div className="column-header">
                  <h2>{columnStatus}</h2>
                  <span className="column-count">{columnTasks.length}</span>
                </div>

                {/* Área de arrastar contém apenas as tarefas */}
                <Droppable droppableId={columnStatus}>
                  {(provided, snapshot) => (
                    <div
                      className={`task-list ${snapshot.isDraggingOver ? 'is-dragging-over' : ''}`}
                      ref={provided.innerRef}
                      {...provided.droppableProps}
                    >
                      {columnTasks.map((task, index) => (
                        <Draggable key={task.id} draggableId={task.id} index={index}>
                          {(provided, snapshot) => (
                            <div
                              className={`task ${snapshot.isDragging ? 'is-dragging' : ''}`}
                              ref={provided.innerRef}
                              {...provided.draggableProps}
                              {...provided.dragHandleProps}
                            >
                              <div className="task-header">
                                <span className="task-title">{task.title}</span>
                                <div className="task-actions">
                                  <button type="button" title="Editar" onClick={() => setFormData(task)}><Edit2 size={16} /></button>
                                  <button type="button" title="Excluir" onClick={() => handleDelete(task.id)}><Trash2 size={16} /></button>
                                </div>
                              </div>
                              {task.description && <div className="task-desc">{task.description}</div>}
                            </div>
                          )}
                        </Draggable>
                      ))}

                      {/* Estado vazio */}
                      {columnTasks.length === 0 && !snapshot.isDraggingOver && (
                        <div className="empty-column">Nenhuma tarefa por aqui</div>
                      )}

                      {provided.placeholder}
                    </div>
                  )}
                </Droppable>

                {/* Botão de adicionar fica fixo no rodapé da coluna */}
                <div className="column-footer">
                  <button
                    className="add-btn"
                    onClick={() => setFormData({ title: '', description: '', status: columnStatus })}
                  >
                    <Plus size={18} style={{ marginRight: '8px' }} />
                    Adicionar Tarefa
                  </button>
                </div>

              </div>
            );
          })}
        </div>
      </DragDropContext>
    </div>
  );
}

export default App;