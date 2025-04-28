import React from 'react';
import styles from './TaskFilters.module.css';

function TaskFilters({ filters, onFilterChange, onResetFilters }) {
    const handleInputChange = (e) => {
        const { name, value } = e.target;
        onFilterChange({ ...filters, [name]: value });
    };

    const handleApplyFilters = () => {
        // Можно добавить дополнительную логику перед применением фильтров
        // В данном случае фильтры уже применяются при изменении инпутов
    };

    return (
        <div className={styles.filters}>
            <div className={styles.filterRow}>
                <div className={styles.filterGroup}>
                    <label htmlFor="company">Компания</label>
                    <select
                        id="company"
                        name="company"
                        value={filters.company}
                        onChange={handleInputChange}
                    >
                        <option value="">Все компании</option>
                        <option value="TechLab">TechLab</option>
                        <option value="DigiCraft">DigiCraft</option>
                        <option value="InnovaTech">InnovaTech</option>
                        <option value="Smart Solutions">Smart Solutions</option>
                    </select>
                </div>

                <div className={styles.filterGroup}>
                    <label htmlFor="category">Категория</label>
                    <select
                        id="category"
                        name="category"
                        value={filters.category}
                        onChange={handleInputChange}
                    >
                        <option value="">Все категории</option>
                        <option value="development">Разработка</option>
                        <option value="design">Дизайн</option>
                        <option value="marketing">Маркетинг</option>
                        <option value="analytics">Аналитика</option>
                        <option value="management">Управление</option>
                    </select>
                </div>

                <div className={styles.filterGroup}>
                    <label htmlFor="status">Статус</label>
                    <select
                        id="status"
                        name="status"
                        value={filters.status}
                        onChange={handleInputChange}
                    >
                        <option value="">Все статусы</option>
                        <option value="new">Новые</option>
                        <option value="in-progress">В работе</option>
                        <option value="completed">Завершенные</option>
                        <option value="urgent">Срочные</option>
                    </select>
                </div>
            </div>

            <div className={styles.filterRow}>
                <div className={styles.filterGroup}>
                    <label htmlFor="search">Поиск</label>
                    <input
                        type="text"
                        id="search"
                        name="search"
                        value={filters.search}
                        onChange={handleInputChange}
                        placeholder="Введите ключевое слово..."
                    />
                </div>

                <div className={styles.filterGroup}>
                    <label htmlFor="deadline">Срок выполнения</label>
                    <select
                        id="deadline"
                        name="deadline"
                        value={filters.deadline}
                        onChange={handleInputChange}
                    >
                        <option value="">Любой срок</option>
                        <option value="today">Сегодня</option>
                        <option value="week">Неделя</option>
                        <option value="month">Месяц</option>
                    </select>
                </div>
            </div>

            <div className={styles.filterActions}>
                <button
                    className={`${styles.btn} ${styles.btnOutline}`}
                    onClick={onResetFilters}
                >
                    Сбросить
                </button>
                <button
                    className={`${styles.btn} ${styles.btnPrimary}`}
                    onClick={handleApplyFilters}
                >
                    Применить
                </button>
            </div>
        </div>
    );
}

export default TaskFilters;