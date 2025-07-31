'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { adminApi, AttributeGroup } from '@/services/admin';
import { toast } from '@/utils/toast';
import GroupForm from './components/GroupForm';
import GroupAttributes from './components/GroupAttributes';

export default function AttributeGroupsPage() {
  const t = useTranslations('admin');
  const [groups, setGroups] = useState<AttributeGroup[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedGroup, setSelectedGroup] = useState<AttributeGroup | null>(
    null
  );
  const [showForm, setShowForm] = useState(false);
  const [showAttributes, setShowAttributes] = useState(false);
  const [isEditing, setIsEditing] = useState(false);

  useEffect(() => {
    loadGroups();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const loadGroups = async () => {
    try {
      setLoading(true);
      const data = await adminApi.attributeGroups.getAll();
      setGroups(data);
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to load groups:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAddGroup = () => {
    setSelectedGroup(null);
    setIsEditing(false);
    setShowForm(true);
    setShowAttributes(false);
  };

  const handleEditGroup = (group: AttributeGroup) => {
    setSelectedGroup(group);
    setIsEditing(true);
    setShowForm(true);
    setShowAttributes(false);
  };

  const handleManageAttributes = (group: AttributeGroup) => {
    setSelectedGroup(group);
    setShowAttributes(true);
    setShowForm(false);
  };

  const handleDeleteGroup = async (group: AttributeGroup) => {
    if (!confirm(t('common.confirmDelete'))) return;

    try {
      await adminApi.attributeGroups.delete(group.id);
      toast.success(t('common.deleteSuccess'));
      await loadGroups();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to delete group:', error);
    }
  };

  const handleSaveGroup = async (data: Partial<AttributeGroup>) => {
    try {
      if (isEditing && selectedGroup) {
        await adminApi.attributeGroups.update(selectedGroup.id, data);
        toast.success(t('common.saveSuccess'));
      } else {
        await adminApi.attributeGroups.create(data);
        toast.success(t('common.saveSuccess'));
      }
      setShowForm(false);
      await loadGroups();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to save group:', error);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">{t('attributeGroups.title')}</h1>
        <button className="btn btn-primary" onClick={handleAddGroup}>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-5 w-5 mr-2"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M12 4v16m8-8H4"
            />
          </svg>
          {t('attributeGroups.addGroup')}
        </button>
      </div>

      <div className="grid grid-cols-1 gap-6">
        <div className="col-span-1">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {groups.length === 0 ? (
              <div className="col-span-full text-center py-8">
                <p className="text-base-content/60">{t('common.noData')}</p>
              </div>
            ) : (
              groups.map((group) => (
                <div key={group.id} className="card bg-base-100 shadow-xl">
                  <div className="card-body">
                    <h2 className="card-title flex items-center gap-2">
                      {group.icon && (
                        <span className="text-2xl">{group.icon}</span>
                      )}
                      {group.display_name}
                    </h2>
                    {group.description && (
                      <p className="text-sm text-base-content/60">
                        {group.description}
                      </p>
                    )}

                    <div className="flex gap-2 items-center mt-2">
                      {group.is_system && (
                        <span className="badge badge-info badge-sm">
                          Системная
                        </span>
                      )}
                      {!group.is_active && (
                        <span className="badge badge-error badge-sm">
                          Неактивна
                        </span>
                      )}
                    </div>

                    <div className="card-actions justify-end mt-4">
                      <button
                        className="btn btn-ghost btn-sm"
                        onClick={() => handleManageAttributes(group)}
                        title={t('attributeGroups.attributes')}
                      >
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          className="h-4 w-4"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"
                          />
                        </svg>
                      </button>
                      <button
                        className="btn btn-ghost btn-sm"
                        onClick={() => handleEditGroup(group)}
                        title={t('common.edit')}
                      >
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          className="h-4 w-4"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                          />
                        </svg>
                      </button>
                      {!group.is_system && (
                        <button
                          className="btn btn-ghost btn-sm text-error"
                          onClick={() => handleDeleteGroup(group)}
                          title={t('common.delete')}
                        >
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            className="h-4 w-4"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                            />
                          </svg>
                        </button>
                      )}
                    </div>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>

      {/* Modal for Group Form */}
      {showForm && (
        <div className="modal modal-open">
          <div className="modal-box w-11/12 max-w-2xl">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-bold">
                {isEditing
                  ? t('attributeGroups.editGroup')
                  : t('attributeGroups.addGroup')}
              </h2>
              <button 
                className="btn btn-sm btn-circle btn-ghost"
                onClick={() => setShowForm(false)}
              >
                ✕
              </button>
            </div>
            <GroupForm
              group={selectedGroup}
              onSave={handleSaveGroup}
              onCancel={() => setShowForm(false)}
            />
          </div>
          <div className="modal-backdrop" onClick={() => setShowForm(false)}></div>
        </div>
      )}

      {/* Modal for Group Attributes */}
      {showAttributes && selectedGroup && (
        <div className="modal modal-open">
          <div className="modal-box w-11/12 max-w-4xl">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-bold">
                {t('attributeGroups.attributes')}: {selectedGroup.display_name}
              </h2>
              <button 
                className="btn btn-sm btn-circle btn-ghost"
                onClick={() => setShowAttributes(false)}
              >
                ✕
              </button>
            </div>
            <GroupAttributes
              group={selectedGroup}
              onUpdate={loadGroups}
              onClose={() => setShowAttributes(false)}
            />
          </div>
          <div className="modal-backdrop" onClick={() => setShowAttributes(false)}></div>
        </div>
      )}
      </div>
    </div>
  );
}
