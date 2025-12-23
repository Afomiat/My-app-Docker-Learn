// Function to fetch and display profile
async function fetchProfile() {
    const emptyState = document.getElementById('emptyProfile');
    const profileContent = document.getElementById('profileContent');
    const statusMsg = document.getElementById('statusMsg');

    try {
        const response = await fetch('/get-profile');
        const data = await response.json();
        
        if (response.ok) {
            // Update display fields
            document.getElementById('displayName').textContent = data.name || '-';
            document.getElementById('displayEmail').textContent = data.email || '-';
            document.getElementById('displayInterests').textContent = data.interests || '-';
            
            // Show content, hide empty state
            emptyState.style.display = 'none';
            profileContent.style.display = 'block';
            
            // Pre-fill form if empty
            if(!document.getElementById('name').value) {
                document.getElementById('name').value = data.name || '';
                document.getElementById('email').value = data.email || '';
                document.getElementById('interests').value = data.interests || '';
            }
        } else {
            // If user not found (404), show empty state
            if (response.status === 404) {
                emptyState.style.display = 'block';
                profileContent.style.display = 'none';
            }
        }
    } catch (error) {
        console.error('Error fetching profile:', error);
    }
}

// Handle Profile Update
document.getElementById('profileForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const btn = this.querySelector('button');
    const originalText = btn.textContent;
    const statusMsg = document.getElementById('statusMsg');
    
    // UI Loading State
    btn.textContent = 'Updating...';
    btn.disabled = true;
    statusMsg.style.display = 'none';
    
    const formData = {
        name: document.getElementById('name').value,
        email: document.getElementById('email').value,
        interests: document.getElementById('interests').value
    };
    
    try {
        const response = await fetch('/update-profile', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(formData)
        });
        
        const data = await response.json();
        
        if (response.ok) {
            // Show success message
            statusMsg.className = 'status-msg success';
            statusMsg.textContent = '✅ Profile updated successfully!';
            statusMsg.style.display = 'block';
            
            // Update display fields directly from form data
            document.getElementById('displayName').textContent = formData.name || '-';
            document.getElementById('displayEmail').textContent = formData.email || '-';
            document.getElementById('displayInterests').textContent = formData.interests || '-';
            
            // Show content, hide empty state
            document.getElementById('emptyProfile').style.display = 'none';
            document.getElementById('profileContent').style.display = 'block';
        } else {
            throw new Error(data.error || 'Update failed');
        }
    } catch (error) {
        // Show error message
        statusMsg.className = 'status-msg error';
        statusMsg.textContent = `❌ ${error.message}`;
        statusMsg.style.display = 'block';
    } finally {
        // Restore button
        btn.textContent = originalText;
        btn.disabled = false;
    }
});

// Refresh button handler
document.getElementById('refreshBtn').addEventListener('click', function() {
    const btn = this;
    const originalText = btn.textContent;
    btn.textContent = 'Refreshing...';
    btn.disabled = true;
    
    fetchProfile().finally(() => {
        btn.textContent = originalText;
        btn.disabled = false;
    });
});

// Load profile on page load
document.addEventListener('DOMContentLoaded', fetchProfile);
